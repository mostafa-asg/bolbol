# bolbol
Are you learning a new language? English? Deutsch? Dutch? Spanish? You learn new words and the next day forget most of them?
The best way to learn new words, is by using them in the sentences and of course repeat them reqularly to move these words from *short-term* memory to *long-term* memory. **bolbol** helps you learn new words by remembering and using them in the full sentences.

## How it works
It is very simple. All you need is a text file to add your sentences to it.

### The text file structure
Between each sentence must exists two `\n`
```
sentence 1\n
\n
sentence 2\n
\n
.
.
.
sentence N
```
### The sentence structure
The word that you want to learn, should be between two **stars** like this:
```
I *go* to school.
```
When you run the application the result should be like this:  
![Sample 1](https://github.com/mostafa-asg/bolbol/blob/master/images/1.png)  
You can star as many word as you want. In this case **bolbol** chooses one of them randomly:
```
Children must be *taught* to *distinguish* between right and wrong.
```
Result:  
![Sample 2](https://github.com/mostafa-asg/bolbol/blob/master/images/2.png)  
In the [languages](https://github.com/mostafa-asg/bolbol/tree/master/languages) folder, you can find these text files.

## How to build
If you do not have [golang](https://golang.org) please install it. Then:
```
git clone https://github.com/mostafa-asg/bolbol/
cd bolbol/
go get github.com/logrusorgru/aurora
go build
```

## How to run
```
./bolbol <YOUR_TEXT_FILE_PATH>
```
## About repository name
In [persian language](https://en.wikipedia.org/wiki/Persian_language), **bolbol** means **nightingale** and there is an idiom, when someone speeks a language very well and fluently:
```
مثل بلبل حرف میزنه
```
Which means he or she speeks like a nightingale.
