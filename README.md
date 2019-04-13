# Kubernetes 구성
선택한 CRI Engine : Docker
선택한 CNI Plugin : Flanner

1. master node에 apt-get install ansible 수행
2. worker node에 apt-get install python 수행
(python이 설치되어 있어야 ansible 작업이 가능)
3. master node의 /etc/ansible/hosts에 다음과 같이 작성

[master]
192.168.0.2(master node ip 혹은 master node domain name) <br/>
[worker]
192.168.0.3(worker node ip 혹은 worker node domain name)

4. 만일 ansible에서 작업을 할 host에 public key로 접근하는 것이 아니라면 추가적인 flag를 적어주어야함

[master]
192.168.0.2 ansible_connection=ssh ansible_ssh_user={user} ansible_ssh_pass={password} <br/>
[worker]
192.168.0.3 ansible_connection=ssh ansible_ssh_user={user} ansible_ssh_pass={password}

5. master node 설정 <br/>
ansible-playbook master-playbook.yml --key-file "awskube.pem" 수행
(aws에서는 public key로 ssh 접속을 하기에 aws 인스턴스에 ansible로 설정한다면 --key-file 옵션 필요)

6. worker node 설정 <br/>
ansible-playbook worker-playbook.yml --key-file "awskube.pem" 수행

7. 성공적으로 작업이 끝나면 kubectl get nodes 명령어로 node 리스트가 정상적으로 뜨는지 확인

# AWS 인스턴스에서 작업 확인하기
1. awskube.pem 공개키를 이용하여 ssh -i awskube.pem ubuntu@ec2-18-219-236-79.us-east-2.compute.amazonaws.com 접속
2. sudo su로 root 계정 전환
3. kubectl get nodes, kubectl get crds로 확인
4. go CRD 프로그램 실행을 위해 환경변수 설정 -> export KUBECONFIG=$HOME/.kube/config
5. /home/ubuntu 경로에 src 바이너리로 go 소스코드 실행 (command : ./src)
