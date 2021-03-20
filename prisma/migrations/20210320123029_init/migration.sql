-- CreateTable
CREATE TABLE "user" (
    "id" TEXT NOT NULL,
    "telegramId" INTEGER NOT NULL,
    "name" TEXT,

    PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "month" (
    "userId" TEXT,
    "id" TEXT NOT NULL,
    "fullDate" TEXT NOT NULL,
    "name" TEXT NOT NULL,

    PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "day" (
    "monthId" TEXT,
    "id" TEXT NOT NULL,
    "date" TEXT NOT NULL,
    "come" TIMESTAMP(3) NOT NULL,
    "go" TIMESTAMP(3) NOT NULL,
    "workHours" INTEGER NOT NULL,

    PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "user.telegramId_unique" ON "user"("telegramId");

-- CreateIndex
CREATE UNIQUE INDEX "month.fullDate_unique" ON "month"("fullDate");

-- CreateIndex
CREATE UNIQUE INDEX "day.date_unique" ON "day"("date");

-- AddForeignKey
ALTER TABLE "month" ADD FOREIGN KEY ("userId") REFERENCES "user"("id") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "day" ADD FOREIGN KEY ("monthId") REFERENCES "month"("id") ON DELETE SET NULL ON UPDATE CASCADE;
