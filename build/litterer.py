import os
import sys
import time
import random


def file_writer(file_path: str, file_size: int):
    with open(file_path, 'w') as file:
        chars_num = int(file_size)
        file.write(str(random.randbytes(chars_num)))


def make_files(files_num: int, file_size: int):
    dir_name = "/data/litter"
    if not os.path.exists(dir_name):
        print("creating dir", dir_name, flush=True)
        os.makedirs(dir_name, exist_ok=True)

    for i in range(0, files_num):
        file_path = f"{dir_name}/file{i}"
        file_writer(file_path, file_size)


if __name__ == '__main__':
    space_to_litter: int = int(sys.argv[1])

    file_size: int = min(space_to_litter, 1000 * 1000) # 10MB per file

    files_count: int = space_to_litter // file_size  # 10MB per file

    make_files(files_count, file_size)

    print("done creating files", flush=True)
    while True:
        time.sleep(1000)
