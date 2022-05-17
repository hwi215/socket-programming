#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <arpa/inet.h>
#include <sys/socket.h>
#include <sys/select.h>
#include <signal.h>
#include <ctype.h>
#include <time.h>
#include <pthread.h>

#define BUF_SIZE 1024
#define PORT 30189

int clntNum;
int serv_sock, clnt_sock;

void error_handling(char *buf);
void (*ctrl_handler)(int);
void ToUp(char*);
void itoa(int n,char s[]);
void reverse(char s[]);


void *timerThread(){
    printf("Number of connected clients = %d\n\n",clntNum);
    while(1){
        sleep(60);
        printf("Number of connected clients = %d\n\n",clntNum);
    }
}
void signalHandler(int sig){
    printf("\nbye bye~\n");
    close(serv_sock);
    exit(1);
}
void ToUp(char *p){
    while(*p){
        *p=toupper(*p);
        p++;
    }
}
void itoa(int n, char s[])
 {
     int i, sign;
 
     if ((sign = n) < 0)  /* record sign */
         n = -n;          /* make n positive */
     i = 0;
     do {       /* generate digits in reverse order */
         s[i++] = n % 10 + '0';   /* get next digit */
     } while ((n /= 10) > 0);     /* delete it */
     if (sign < 0)
         s[i++] = '-';
     s[i] = '\0';
     reverse(s);
 }
  /* reverse:  reverse string s in place */
 void reverse(char s[])
 {
     int i, j;
     char c;
 
     for (i = 0, j = strlen(s)-1; i<j; i++, j--) {
         c = s[i];
         s[i] = s[j];
         s[j] = c;
     }
 }

int main(){

    //timer thread
    pthread_t timer;
    int threadErr;

    if(threadErr = pthread_create(&timer,NULL,timerThread,NULL)){
        //thread err
        printf("Thread Err = %d",threadErr);
    }

    //바이트순서로 정렬된 정수형IP주소를 우리눈으로 쉽게 인식할 수 있는 문자열 형태로 반환
    char *clnt_Ip; 

    char index='5';

    int count[100]={0,};
    struct timespec begin,end[100];
    int start[100]={0,};

    struct sockaddr_in serv_adr, clnt_adr;

    struct timeval timeout;
    fd_set reads, cpy_reads;

    socklen_t adr_sz;
    int fd_max, str_len, fd_num, i;

    serv_sock=socket(PF_INET,SOCK_STREAM,0);
    memset(&serv_adr,0,sizeof(serv_adr));
    serv_adr.sin_family=AF_INET;
    serv_adr.sin_addr.s_addr=htonl(INADDR_ANY);
    serv_adr.sin_port=htons(PORT);

    if(bind(serv_sock,(struct sockaddr*) &serv_adr,sizeof(serv_adr))==-1)
        error_handling("bind() error 1분뒤에 다시 연결 or 다른 port번호로 접속해주세요");
    if(listen(serv_sock,5)==-1)
        error_handling("listen() error");

    printf("Server is ready to receive on port %d\n\n",ntohs(serv_adr.sin_port));

    FD_ZERO(&reads);
    FD_SET(serv_sock,&reads);
    fd_max=serv_sock;

    //ctrl+c
    ctrl_handler=signal(SIGINT,signalHandler);

    while(1){

        cpy_reads=reads;
        timeout.tv_sec=5;
        timeout.tv_usec=5000;

        if((fd_num=select(fd_max+1,&cpy_reads,0,0,&timeout))==-1)
            break;
        if(fd_num==0)
            continue;
        
        
        for(i=0;i<fd_max+1;i++){
            if(FD_ISSET(i,&cpy_reads))
            {
                int startTime;
                if(i==serv_sock)  //connection request!
                {
                    clntNum++;

                    adr_sz=sizeof(clnt_adr);
                    clnt_sock=accept(serv_sock,(struct  sockaddr*)&clnt_adr,&adr_sz);
                    printf("Client %d connected. Number of connected clients %d \n\n",clnt_sock-3,clntNum);

                    //Q4
                    clock_gettime(CLOCK_MONOTONIC,&begin);
                    start[clnt_sock-3]=(int)begin.tv_sec;

                    FD_SET(clnt_sock,&reads);
                    if(fd_max<clnt_sock)
                        fd_max=clnt_sock;
                    
                }
                else        //read message!
                {
                    //Q3
                    count[i-3]++;

                    char buf[BUF_SIZE];
                    clnt_Ip=inet_ntoa(clnt_adr.sin_addr);

                    //Q2
                    char clnt_address[BUF_SIZE]={0,};
                    sprintf(clnt_address,"%s:%d",clnt_Ip,ntohs(clnt_adr.sin_port));

                    printf("Connection request from %s:%d \n",clnt_Ip,ntohs(clnt_adr.sin_port));
                    str_len = read(i,buf,BUF_SIZE-1);
                    
                    buf[str_len]=0;
                    index=buf[0];

                    if(index=='1'){
                        char *low=strtok(buf,"1");
                        ToUp(low);
                        write(i,low,str_len);
                    }
                    else if(index=='2'){
                        write(i,clnt_address,BUF_SIZE);
                    }
                    else if(index=='3'){
                        char countNum[BUF_SIZE]={0,};
                        itoa(count[i-3],countNum);
                        write(i,countNum,BUF_SIZE);
                    }
                    else if(index=='4'){
                        clock_gettime(CLOCK_MONOTONIC,&end[i-3]);
                        int endtime = (int)end[i-3].tv_sec;
                        int second = endtime-start[i-3];
                        int minute = second/60;
                        int hours = minute/60;
                        char runtime[BUF_SIZE]={0,};
                        sprintf(runtime,"%02d:%02d:%02d",hours,minute,second);
                        write(i,runtime,BUF_SIZE);

                    }else if(index=='5'){
                        clntNum--;
                        FD_CLR(i,&reads);
                        printf("Client %d disconnected. Number of connected clients = %d\n",i-3,clntNum);
                    }
                }
            }
        }
    }
    close(serv_sock);
    return 0;
}

void error_handling(char *buf)
{
    fputs(buf,stderr);
    fputc('\n',stderr);
    exit(1);
}