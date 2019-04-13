# KubernetesCRD

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
