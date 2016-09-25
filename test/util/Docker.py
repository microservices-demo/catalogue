import re
from subprocess import Popen, PIPE
from random import random
from time import sleep

# From http://blog.bordage.pro/avoid-docker-py/


class Docker:
    def kill_and_remove(self, ctr_name):
        command = ['docker', 'rm', '-f', ctr_name]
        try:
            self.execute(command)
        except RuntimeError as e:
            print(e)
            pass

    def ensure_output(self, container_name, match_string, limit=30):
        command = ['docker', 'logs', container_name]
        matched = False
        while not matched:
            out = Docker().execute(command)
            matched = True if match_string in out else False
            if matched:
                return True
            if limit == 0:
                return False
            else:
                sleep(1)
                limit = limit - 1
 
    def random_container_name(self, prefix):
        retstr = prefix + '-'
        for i in range(5):
            retstr += chr(int(round(random() * (122-97) + 97)))
        return retstr

    def get_container_ip(self, ctr_name):
        command = ['docker', 'inspect',
                   '--format', '\'{{.NetworkSettings.IPAddress}}\'',
                   ctr_name]
        return re.sub(r'[^0-9.]*', '', self.execute(command))

    def execute(self, command):
        p = Popen(command, stdout=PIPE, stderr=PIPE)
        out = p.stdout.read()
        stderr = p.stderr.read()
        if p.wait() != 0:
            p.stdout.close()
            p.stderr.close()
            raise RuntimeError(str(stderr, 'utf-8'))
        p.stdout.close()
        p.stderr.close()
        return str(out, 'utf-8')

    def start_container(self, container_name="", image="", cmd="", host="", env=[]):
        command = ['docker', 'run', '-d']
        if image == "":
            raise RuntimeError(str("Image can't be empty", 'utf-8'))
        
        command.extend(["--name", container_name]) if container_name != "" else 0
        command.extend(["-h", host]) if host != "" else 0

        if env != "":
            [command.extend(["--env", "{}={}".format(x[0], x[1])]) for x in env]
            
        command.append(image)
        
        if cmd != "":
            command.append(cmd)
        
        self.execute(command)
