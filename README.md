# Kubernetes 구성
선택한 CRI Engine : Docker <br/>
선택한 CNI Plugin : Flannel <br/>

1. master node에 apt-get install ansible 작업을 수행합니다.
2. worker node에 apt-get install python 작업을 수행합니다.
(python이 설치되어 있어야 ansible 작업이 가능합니다)
3. master node의 /etc/ansible/hosts에 아래의 예시를 참고하여 작성합니다

[master] <br/>
192.168.0.2(master node ip 혹은 master node domain name) <br/>
[worker] <br/>
192.168.0.3(worker node ip 혹은 worker node domain name) <br/>

4. 만일 ansible에서 작업을 할 host에 public key로 접근하는 것이 아니라면 추가적인 flag를 적어주어야합니다

[master] <br/>
192.168.0.2 ansible_connection=ssh ansible_ssh_user={user} ansible_ssh_pass={password} <br/>
[worker] <br/>
192.168.0.3 ansible_connection=ssh ansible_ssh_user={user} ansible_ssh_pass={password} <br/>

5. master node 설정 <br/>
ansible-playbook master-playbook.yml --key-file "awskube.pem" 명령어를 수행합니다
(aws에서는 public key로 ssh 접속을 하기에 aws 인스턴스에 ansible로 설정한다면 --key-file 옵션이 필요합니다)
(접속에 필요한 public key file은 key 디렉토리에 있습니다)

6. worker node 설정 <br/>
ansible-playbook worker-playbook.yml --key-file "awskube.pem" 명령어를 수행합니다

7. 성공적으로 작업이 끝나면 sudo su 명령으로 root 계정으로 전환 후 kubectl get nodes 명령어로 node 리스트가 정상적으로 뜨는지 확인합니다

# AWS 인스턴스에서 작업 확인하기
1. awskube.pem 공개키를 이용하여 ssh -i awskube.pem ubuntu@ec2-18-219-236-79.us-east-2.compute.amazonaws.com에 접속합니다
2. sudo su로 root 계정으로 전환합니다
3. kubectl get nodes, kubectl get crds로 셋팅을 확인합니다
4. go CRD 프로그램 실행을 위해 환경변수를 설정합니다 -> export KUBECONFIG=$HOME/.kube/config
5. /home/ubuntu 경로에 src 바이너리로 go 소스코드 실행합니다 (command : ./src)
6. 동일한 디렉토리에 crd.log 로그가 생성된 것을 확인할 수 있습니다
