import React, { useState } from "react";
import { Box, Typography, Stepper, Step, StepLabel, Paper } from "@mui/material";

import AddressForm from "./Checkout/AddressForm";
import PaymentForm from "./Checkout/PaymentForm";

const steps = ["Shipping address", "Payment details"]

const Confirmation = () => {
    return (
        <Typography variant="h5">
            Confirmation
        </Typography>
    )
}

const Form = ({step, cart}) => {
    return (
        step.state === 0 ? <AddressForm next={step.next}/> : <PaymentForm cart={cart} prev={step.prev}/>
    )
}

const Payment = ({cart}) => {
    const [activeStep, setActiveStep] = useState(0)

    const nextStep = () => setActiveStep((prevActiveStep) => prevActiveStep + 1)
    const prevStep = () => setActiveStep((prevActiveStep) => prevActiveStep - 1)

    return (
        <Box sx={{
            display: "flex",
            justifyContent: "center",
            marginTop: 5,
        }}>
            <Paper sx={{
                width: 1/3}}>
                <Typography variant="h5" align="center" gutterBottom>
                    Checkout
                </Typography>
                <Stepper activeStep={activeStep}>
                    {steps.map((step) => (
                        <Step key={step}>
                            <StepLabel>{step}</StepLabel>
                        </Step>
                    ))}
                </Stepper>
                {activeStep === steps.length ? <Confirmation/> : <Form cart={cart} step={{state: activeStep, next: nextStep, prev: prevStep}}/>}
            </Paper>
        </Box>
    )
}

export default Payment