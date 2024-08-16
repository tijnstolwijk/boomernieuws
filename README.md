# Boomernieuws

Read teletekst from the terminal.
The code is horrible, please use any other teletekst reader,
except of course the electron client xD jkjk.
BTW: just build it yourself, i won't update the releases.

## Usage

The vim keys h, j, k, l are used to navigate pages and subpages.

To visit a certain page press : (colon) then the address of the page.
For example, to visit the second subpage of page 711:
```
:711-2
```

To quit boomernieuws simply press:
```
:q
```

## Downloading a single page

To download a single page use the --output flag.
For example, to download page 101 to example.txt, use this command:
```
boomernieuws --output example.txt 101
```

## Printing a single page (non-interactive mode)

To print page 101, enter this command:
```
boomernieuws 101
```

## Archiveringsdienst

Dutch for archiving service.
This will download all teletekst pages every hour, 
but will only store it if it is actually different from the 
other most recent version of that page. 
It does this by searching through the timestamps from new to old
until it finds the page, then compares it to the downloaded page.
If it's different the page will be saved in a new timestamp.
This results in a system that will only ever save unique pages.

"De archiveringsdienst" can also read a certain teletekst page 
from any given moment in time (this time must be between the first recording and the present).

To get the most recent recorded version of page 101 you run this command:
```
$ archiveringsdienst read latest 101
```

To get page 101 from timestamp 1723754183 you run this command:
```
$ archiveringsdienst read 1723754183 101
```

Archiveringsdienst will look to the most recent recorded timestamp for the given timestamp
if the page is not present in this timestamp it looks to the timestamps before it
until it finds it.
