import os
import sys


def file_writer(file_path, file_size):
    with open(file_path, 'w') as file:
        chars_num = int(file_size)
        file.write("0" * chars_num)


def make_files(dir_name, files_num, file_size):
    dir_name = '/data/' + dir_name
    if not os.path.exists(dir_name):
        os.makedirs(dir_name, exist_ok=True)

    for i in range(0, files_num):
        file_path = dir_name + '/file' + str(i)
        file_writer(file_path, file_size)


if __name__ == '__main__':
    arg_list = sys.argv

    dir_name = str(arg_list[1])
    files_num = int(arg_list[2])
    file_size = int(arg_list[3])
    
    make_files(dir_name, files_num, file_size)
