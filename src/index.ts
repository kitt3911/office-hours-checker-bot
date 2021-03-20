import { PrismaClient } from '@prisma/client'
import { Telegraf } from 'telegraf'
import { botToken } from './config'
import { Day } from './interfaces/day.interface'
import { createDay, findOneDay } from './service/day.service'
import { validationMonth, validationUser } from './service/validation.service'
import { formatHours } from './utils/formatDate'

const prisma = new PrismaClient()

async function main() {
    const bot = new Telegraf(botToken)
    bot.start(async (ctx) => {
        const user = {
            name: ctx.from.username,
            telegramId: ctx.from.id
        }
        const validateUser = await validationUser(prisma,user)

        ctx.reply(`Welcome ${validateUser.name}`)
    })

    bot.hears(/(\d{2})\/((\d{2})|(\d{1}))\s(\d)/, async (ctx) => {
        const setDay = ctx.match[1] + '/' + ctx.match[2]
        const hours = ctx.match[5]
        console.log(ctx.match)
        const user = {
            name: ctx.from.username,
            telegramId: ctx.from.id
        }
        const validateUser = await validationUser(prisma,user)
        const validateMonth = await validationMonth(prisma,setDay,validateUser.id)
      //  console.log(validateMonth)
        const newDay = await createDay(prisma,validateMonth.id,Number(hours),setDay)
        ctx.reply(JSON.stringify(newDay))
    })

    bot.hears(/(\d{2})\/((\d{2})|(\d{1}))/,async (ctx)=> {
        const setDay = ctx.match[1] + '/' + ctx.match[2]
        const user = {
            name: ctx.from.username,
            telegramId: ctx.from.id
        }
        const validateUser = await validationUser(prisma,user)
        const validateMonth = await validationMonth(prisma,setDay,validateUser.id)
        const day = await findOneDay(prisma,validateMonth.id) 
    
        if(day){
            const formatDate = formatHours(day.workHours)
            ctx.reply(`
                date: ${day.date} /n,
                ${formatDate}
            `)
        }

    })

    bot.launch()
}

main().catch(e => {
    throw e
}).finally(async () => {
    await prisma.$disconnect()
})