# Kubernetes 구성
1. master node에 apt-get install ansible 작업을 수행합니다.
2. worker node에 apt-get install python 작업을 수행합니다.
(python이 설치되어 있어야 ansible 작업이 가능합니다)
3. master node의 /etc/ansible/hosts에 아래의 예시를 참고하여 작성합니다

[master] <br/>
192.168.0.2(master node ip 혹은 master node domain name) <br/>
[worker] <br/>
192.168.0.3(worker node ip 혹은 worker node domain name) <br/>

4. 만일 ansible에서 작업을 할 host에 public key로 접근하는 것이 아니라면 추가적인 option을 적어주어야합니다
[master] <br/>
192.168.0.2 ansible_connection=ssh ansible_ssh_user={user} ansible_ssh_pass={password} <br/>
[worker] <br/>
192.168.0.3 ansible_connection=ssh ansible_ssh_user={user} ansible_ssh_pass={password} <br/>

5. master node 설정 <br/>
ansible-playbook master-playbook.yml <br/>

6. worker node 설정 <br/>
ansible-playbook worker-playbook.yml <br/>

# CRD 소스코드 실행
1. sudo su로 root 계정으로 전환합니다
2. kubectl get nodes, kubectl get crds로 셋팅을 확인합니다
3. go CRD 프로그램 실행을 위해 환경변수를 설정합니다 -> export KUBECONFIG=$HOME/.kube/config
4. /home/ubuntu 경로에 src 바이너리로 go 소스코드 실행합니다 (command : ./src) <br/>
5. 만일 src 파일이 실행권한이 없다면 chmod +x src 명령어를 수행한 후 ./src 해주세요.
6. 동일한 디렉토리에 crd.log 로그가 생성된 것을 확인할 수 있습니다
