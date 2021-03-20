import { Month } from "@prisma/client";

export interface User {
    id?: string
    name?: string | null
    telegramId: number
    months?: Month[]
}