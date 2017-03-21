import os
import unittest
from os.path import expanduser

from util.Docker import Docker


class GoServices(unittest.TestCase):
    def test_go(self):
        script_dir = os.path.dirname(os.path.realpath(__file__))
        code_dir = script_dir + "/.."
        home = expanduser("~")
        goPath = os.environ['GOPATH']
        command = ['docker', 'run', '--rm', '-v', goPath + ':/go/src/', '-v', code_dir + ':/go/src/github.com/microservices-demo/catalogue', '-w', '/go/src/github.com/microservices-demo/catalogue', '-e', 'GOPATH=/go/', 'golang:1.7', 'go', 'test', '-v', '-covermode=count', '-coverprofile=coverage.out']

        print(Docker().execute(command))

if __name__ == '__main__':
    unittest.main()
