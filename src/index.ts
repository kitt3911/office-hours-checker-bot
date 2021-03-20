import { PrismaClient } from '@prisma/client'
import {Telegraf} from 'telegraf'
import { botToken } from './config'


const prisma = new PrismaClient()

async function main() {
    const bot = new Telegraf(botToken)
    bot.start((ctx) => ctx.reply('Welcome'))
    bot.launch()
}

main().catch(e => {
    throw e
}).finally(async ()=> {
    await prisma.$disconnect()
})