#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <signal.h>
#include <unistd.h>
#include <arpa/inet.h>
#include <sys/socket.h>
#include <netdb.h>
#include <time.h>
#include <ctype.h>

#define BUF_SIZE 1024
#define PORT 25845
#define DOMAIN "nsl2.cau.ac.kr"

int sock;
clock_t start,end;

void error_handling(char *message);
void (*ctrl_handler)(int);
void upper(int sock);
void getAddress(int sock);
void getrequest(int sock);
void getRuntime(int sock);

void error_handling(char *message)
{
    fputs(message,stderr);
    fputc('\n',stderr);
    exit(1);
}

//ctrl+c
void signalHandler(int sig){
    write(sock,"5",BUF_SIZE);
    printf("\nbye bye~\n");
    close(sock);
    exit(1);
}

void upper(int sock){
    char lower[BUF_SIZE-1];
    char message[BUF_SIZE]="1";
    printf("Input lowercase sentence: ");
    scanf(" %s", lower); 
    fflush(stdin);
    //fgets(lower,BUF_SIZE-1,stdin);
    strcat(message,lower);
    start = clock();
    write(sock,message,strlen(message)-1);
    int str_len=read(sock,message,BUF_SIZE);
    end = clock();
    double rtt = (double)end-start;
    message[str_len]=0;
    printf("Reply from server: %s \n",message);
    printf("RTT : %.3f ms\n",rtt);
}

void getAddress(int sock){
    char addr[BUF_SIZE]={0,};
    start = clock();
    write(sock,"2",BUF_SIZE-1);
    read(sock,addr,BUF_SIZE);
    end = clock();
    fflush(stdin);
    char *ip = strtok(addr,":");
    char *port=strtok(NULL,":");
    double rtt = (double)end-start;
    printf("Reply from server: client IP = %s , port = %s\n",ip,port);
    printf("RTT : %.3f ms\n",rtt);
}
void getreguest(int sock){
    start = clock();
    write(sock,"3",BUF_SIZE-1);
    char request[BUF_SIZE]={0,};
    read(sock,request,BUF_SIZE);
    end = clock();
    fflush(stdin);
    double rtt = (double)end-start;
    printf("Reply from server: client request number = %s\n",request);
    printf("RTT : %.3f ms\n",rtt);
}
void getRuntime(int sock){
    char runtime[BUF_SIZE]={0,};
    start = clock();
    write(sock,"4",BUF_SIZE-1);
    read(sock,runtime,BUF_SIZE);
    end = clock();
    fflush(stdin);
    double rtt = (double)end-start;
    printf("Reply from server: runtime = %s\n\n",runtime);
    printf("RTT : %.3f ms\n",rtt);
}

int main(){
    int recv_len,recv_cnt;
    struct sockaddr_in serv_adr;

    //domain ->ip
    struct hostent *he;
    if((he = gethostbyname(DOMAIN))==NULL){
        error_handling("no server");
    }

    //ctrl+c
    ctrl_handler=signal(SIGINT,signalHandler);

    sock=socket(PF_INET,SOCK_STREAM,0);
    if(sock==-1)
        error_handling("socket() error");

    memset(&serv_adr,0,sizeof(serv_adr));
    serv_adr.sin_family=AF_INET;
    serv_adr.sin_addr.s_addr=inet_addr(inet_ntoa(*(struct in_addr*)he->h_addr_list[0]));
    serv_adr.sin_port=htons(PORT);

    if(connect(sock,(struct sockaddr*)&serv_adr,sizeof(serv_adr))==-1)
        error_handling("connect() error!");
    
    while (1)
    {
        printf("<Menu>\n");
        printf("1) convert text to UPPER-case\n");
        printf("2) get my IP address and port number\n");
        printf("3) get server request count\n");
		printf("4) get server running time\n");
		printf("5) exit\n");
        printf("Input option : ");
        rewind(stdin);
        char input=0;
        scanf(" %c",&input);
        if(input =='1'){
            upper(sock);
        }else if(input=='2'){
            getAddress(sock);
        }else if(input=='3'){
            getreguest(sock);
        }else if(input=='4'){
            getRuntime(sock);
        }else if(input=='5'){
            printf("Bye bye~\n");
            write(sock,"5",1);
            close(sock);
            exit(0);
        }

    }
    
}
