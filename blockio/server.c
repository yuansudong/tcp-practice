
#include <stdio.h>
#include <sys/socket.h>
#include <sys/types.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <unistd.h> 
#include <string.h> 
#include <stdlib.h> 
 
int main(int argc, char *argv[])
{

	int port = 12345;
	//1 创建tcp通信socket
	int socket_fd = socket(AF_INET, SOCK_STREAM, 0);
	if(socket_fd == -1)
	{
		perror("创建tcp通信socket失败!\n");
		return -1;
	}
	struct sockaddr_in server_addr = {0};
	server_addr.sin_family = AF_INET; 
	server_addr.sin_port = htons(port);
	server_addr.sin_addr.s_addr = INADDR_ANY;
	int ret = bind(socket_fd, (struct sockaddr *) &server_addr, sizeof(server_addr) );
	if(ret == -1)
	{
		perror("bind failed!\n");
		return -1;
	}
	//3 设置监听队列，设置为可以接受5个客户端连接
	ret = listen(socket_fd, 5);
	if(ret == -1)
	{
		perror("listen falied!\n");
		return -1;
	}
	printf("服务启动成功!\n");
	struct sockaddr_in client_addr = {0};//用来存放客户端的地址信息
	int len = sizeof(client_addr);
	int new_socket_fd = -1;//存放与客户端的通信socket
	new_socket_fd = accept( socket_fd, (struct sockaddr *)&client_addr, &len);
	if(new_socket_fd == -1)
	{
		perror("accpet error!\n");
		return -1;
	}
	printf("IP:%s, PORT:%d [connected]\n", inet_ntoa(client_addr.sin_addr), ntohs(client_addr.sin_port));
	char buf[40960] = {0};
	while (1)
	{
	 int n = read(new_socket_fd, buf, sizeof(buf));
	 if (n == 0 ||  n == -1){
		 break;
	 }
	 printf("receive msg:%s\n", buf);//打印消息
	}
	printf("接收完毕！\n");
	//5 关闭socket
	close(new_socket_fd);
	close(socket_fd);
	return 0;
}