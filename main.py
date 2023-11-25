import os
import random
import time
from enum import Enum
from typing import Optional

import uvicorn
from fastapi import FastAPI
from datetime import date

from pydantic import BaseModel

app = FastAPI()


RATES_DATA = {
    "PLN": {
        "USD": 4.3188,
        "EUR": 4.5892,
        "PLN": 1.0,
    }
}


class BaseCurrency(str, Enum):
    PLN = "PLN"


class Currency(str, Enum):
    EUR = "EUR"
    PLN = "PLN"
    USD = "USD"


class HistoricalRateResponse(BaseModel):
    for_date: date
    base_currency: BaseCurrency
    currency: Currency
    rate: Optional[float] = None


@app.get("/api/v1/historical_rates")
async def historical_rates(
    for_date: date, base_currency: BaseCurrency, currency: Currency
):
    time.sleep(random.randrange(200, 300) / 1000)
    if for_date != date(2023, 9, 25) and currency != base_currency:
        return HistoricalRateResponse(
            for_date=for_date,
            base_currency=base_currency,
            currency=currency,
            rate=None,
        )

    rate = RATES_DATA.get(base_currency).get(currency)

    return HistoricalRateResponse(
        for_date=for_date,
        base_currency=base_currency,
        currency=currency,
        rate=rate,
    )


if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000)
