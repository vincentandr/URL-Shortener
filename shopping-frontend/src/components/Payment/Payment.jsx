import React, { useEffect, useState } from "react";
import { useSelector, useDispatch } from "react-redux";
import { Box, Typography, Stepper, Step, StepLabel, useMediaQuery, useTheme, Stack} from "@mui/material";
import { Elements } from "@stripe/react-stripe-js"
import { loadStripe } from "@stripe/stripe-js"

import AddressForm from "./Checkout/AddressForm";
import PaymentForm from "./Checkout/PaymentForm";
import Confirmation from "./Checkout/Confirmation";
import defaultFormData from "../../constants/FormData";
import { fetchDraftOrder } from "../../actions";
import EmptyCart from "./Checkout/EmptyCart";

const steps = ["Shipping address", "Payment details", "Confirmation"]

const stripePromise = loadStripe(process.env.REACT_APP_STRIPE_PUBLIC_KEY)

const getSelectors = (state) => ({
  payment: state.payment,
});

const Form = ({step, payment, formData}) => {
    return (
        (step.state === 0 && <AddressForm next={step.next} formData={formData}/>) 
        || (step.state === 1 && <PaymentForm payment={payment} step={step} formData={formData}/>)
        || (step.state === 2 && <Confirmation payment={payment} />)
    )
}

const PaymentContent = ({payment}) => {
    const [activeStep, setActiveStep] = useState(0)
    const [formData, setFormData] = useState(defaultFormData)

    const nextStep = () => setActiveStep((prevActiveStep) => prevActiveStep + 1)
    const prevStep = () => setActiveStep((prevActiveStep) => prevActiveStep - 1)

    return (
        <Stack spacing={1} sx={{
            height: "100%"
        }}>
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
            {activeStep === steps.length ? <Confirmation/> : 
            <Form 
                payment={payment} 
                formData={{state: formData, set: setFormData}} 
                step={{state: activeStep, next: nextStep, prev: prevStep}}/>}
        </Stack>
    )
}

const Payment = () => {
    const theme = useTheme()
    const smallPhone = useMediaQuery(theme.breakpoints.down("sm"))
    const tablet = useMediaQuery(theme.breakpoints.down("md"))

    const dispatch = useDispatch();
    const { payment } = useSelector(getSelectors);

    useEffect(() => {
        dispatch(fetchDraftOrder("user1"))
    }, [dispatch])

    return (
        (payment.secret_key ? (
            <Elements stripe={stripePromise} options={{
                clientSecret: payment.secret_key
            }}>
                <Box sx={{
                    display: "flex",
                    justifyContent: "center",
                    height: "100vh",
                    pl: 2,
                    pr: 2,
                }}>
                    <Box sx={{
                        width: ((smallPhone && 320) || (tablet && 450) || (!tablet && 600)),
                        height: "100%",
                    }}>
                        <PaymentContent payment={payment}/>
                    </Box>
                </Box>
            </Elements>
        ) : <EmptyCart/>)
    )
}

export default Payment