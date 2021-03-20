import { Day } from "./day.interface";

export interface Month {
    id: string
    userId?: string | null
    fullDate: String 
    name: String
    day?: Day[],
}