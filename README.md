# kpp - a lightweight tool for imposing limits on docker volumes usage

kpp (from Russian: КПП, контрольно-пропускной пункт, meaning "border point") is a lightweight tool dedicated



## litterer
To test kpp, you can use build/litterer.py. 
It instantly fills a volume with random files constituting provided number in size
```bash
docker build -t litterer ./build/
```
To run litterer do
```bash
docker run -v <folder_to_litter>:/data litterer <space_amount>
```