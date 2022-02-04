import React, { useState } from "react";
import { Box, Typography, Stepper, Step, StepLabel, Paper} from "@mui/material";

import AddressForm from "./Checkout/AddressForm";
import PaymentForm from "./Checkout/PaymentForm";
import Confirmation from "./Checkout/Confirmation";

const steps = ["Shipping address", "Payment details", "Confirmation"]

const Form = ({step, cart, formData}) => {
    return (
        (step.state === 0 && <AddressForm next={step.next} formData={formData}/>) 
        || (step.state === 1 && <PaymentForm cart={cart} step={step} formDataState={formData.state}/>)
        || (step.state === 2 && <Confirmation cart={cart} />)
    )
}

const Payment = ({cart}) => {
    const [activeStep, setActiveStep] = useState(0)
    const [formData, setFormData] = useState({
        firstName: "",
        lastName: "",
        address: "",
        email: "",
        area: "",
        postal: "",
    })

    const nextStep = () => setActiveStep((prevActiveStep) => prevActiveStep + 1)
    const prevStep = () => setActiveStep((prevActiveStep) => prevActiveStep - 1)

    return (
        <Box sx={{
            display: "flex",
            justifyContent: "center",
            marginTop: 5,
        }}>
            <Paper sx={{
                width: 1/3,
                pl: 3,
                pr: 3,}}>
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
                    {activeStep === steps.length ? <Confirmation/> : <Form cart={cart} formData={{state: formData, set: setFormData}} step={{state: activeStep, next: nextStep, prev: prevStep}}/>}
            </Paper>
        </Box>
    )
}

export default Payment