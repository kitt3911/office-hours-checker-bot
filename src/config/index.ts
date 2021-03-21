import { PrismaClient } from "@prisma/client"

require('dotenv').config()

export const botToken:string = process.env.BOT_TOKEN || ""

const prisma = new PrismaClient()

export {prisma}