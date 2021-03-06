# office-hours-checker-bot
---------
![illustration](https://news.noteru.com/wp-content/uploads/2020/07/08/boty-min.png)

just a telegram bot to monitor the number of hours worked

## How to use ? 

[![Telegram Badge](https://img.shields.io/badge/-Telegram-0088cc?style=flat-square&logo=Telegram&logoColor=white)](https://web.telegram.org/#/im?p=@officeHoursCheckerBot)

Write a message to the bot **@officeHoursCheckerBot**

### Commands :


- <code><button height="27">/come</button></code> - Mark that I came ( starts counting the time )
- <code><button height="27">/go</button></code> - Mark that I left ( stops counting time )
- <code><button height="27">/show</button></code> - Show hours for the day
- <code><button height="27">/show_all_info</button></code> - show the work schedule for the month
- <code><button height="27">get {mm/dd}</button></code> - show the work schedule for the day {mm/dd}
- <code><button height="27">set {mm/dd} {work hour}</button></code> - set the work schedule for the day {mm/dd} {workhour}

## Example :

```javascript
set 02/01 8
get 02/01
```

Result:
```
date: 2021/2/1 
time: 8:00
```
```javascript
set 02/01 8:30
get 02/01
```

Result:
```
date: 2021/2/1 
time: 8:30
```

## For contributors : 

### How to start : 

1. clone this repo
2. start : 
```bash
$ docker-compose up
$ prisma migrate dev --name init
$ pnpm install
$ pnpm run start
```


### How to participate in the project ?
1. Create your feature branch (```git checkout -b my-new-feature```)
2. Commit your changes (```git commit -am 'Add some feature'```)
3. Push to the branch (```git push origin my-new-feature```)
4. Create new Pull Request
5. ???
6. PROFIT!!!
