import time
import sys

if len(sys.argv) > 1:
    server = sys.argv[2]
    api_key = sys.argv[4]
else:
    server = input("server > ")
    api_key = input("username > ")

print("-"*10)
print("starting execution")
print("server:", server)
print("api_key:", api_key)

time.sleep(0.1)

print("stopping execution")

