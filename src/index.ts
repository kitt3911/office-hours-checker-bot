import { PrismaClient } from '@prisma/client'
import { Telegraf } from 'telegraf'
import { botToken } from './config'
import moment from 'moment'

const prisma = new PrismaClient()

async function main() {
    const bot = new Telegraf(botToken)
    bot.start(async (ctx) => {
        const user = {
            name: ctx.from.username,
            telegramId: ctx.from.id
        }
        let findUser = await prisma.user.findFirst({
            where: { telegramId: user.telegramId }
        })
        if (!findUser) {
            findUser = await prisma.user.create({
                data: { ...user }
            })
        }

        let thisDate = moment(new Date(), 'DD/MM/YYYY')
        let month = thisDate.format('M')
        let year = thisDate.format('YYYY')
        const thisFormatDate = year + '/' + month
        const findMonth = await prisma.month.findFirst({
            where: {
                fullDate: thisFormatDate,
                user: findUser,
                userId: findUser.id
            }
        })

        if (!findMonth){
           const modelMonth = await prisma.month.create({
               data:{
                   fullDate: thisFormatDate,
                   userId: findUser.id,
                   name: month
               }
           })
        }

        ctx.reply(`Welcome ${findUser.name}`)
    })

    bot.command("come", (ctx) => {

    })


    bot.launch()
}

main().catch(e => {
    throw e
}).finally(async () => {
    await prisma.$disconnect()
})