import React from "react"
import { Button } from "@mui/material"
import Review from "./Review"

const PaymentForm = ({cart}) => {
    return (
        <>
            <Review cart={cart}/>
        </>
    )
}

export default PaymentForm