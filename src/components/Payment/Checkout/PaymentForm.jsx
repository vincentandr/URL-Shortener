import React, { useEffect, useState } from "react"
import { useDispatch } from "react-redux"
import { Alert, Box, Typography, CircularProgress } from "@mui/material"
import { PaymentElement, useElements, useStripe } from "@stripe/react-stripe-js"

import Review from "./Review"
import { formatCurrency } from "../../../helpers/Utils"
import { makePayment } from "../../../actions"
import { BlackButton } from "../../../theme"

const PaymentForm = ({payment, step, formData}) => {
    const dispatch = useDispatch()
    const elements = useElements()
    const stripe = useStripe()
    const [message, setMessage] = useState("")
    const [isLoading, setLoading] = useState(false)

    useEffect(() => {
        if (!stripe || !payment.secret_key) {
            return;
        }

        stripe.retrievePaymentIntent(payment.secret_key).then(({ paymentIntent }) => {
            switch (paymentIntent.status) {
                case "succeeded":
                    setMessage("Payment succeeded!");
                    break;
                case "processing":
                    setMessage("Your payment is processing.");
                    break;
                case "requires_payment_method":
                    setMessage("Your payment was not successful, please try again.");
                    break;
                default:
                    setMessage("Something went wrong.");
                    break;
            }
        });
    }, [stripe]);

    const handleSubmit = async (event) => {
        event.preventDefault();

        if (!stripe || !elements) return;

        setLoading(true)

        const { error } = await stripe.confirmPayment({
            elements,
            redirect: "if_required",
        });

        setLoading(false)
        
        if (error) {
            if (error.type === "card_error" || error.type === "validation_error") {
                setMessage(error.message);
            } else {
                setMessage("An unexpected error occured.");
            }
        } else {
            const orderData = {
                order_id: payment.order.order_id,
                customer: {
                    first_name: formData.state.first_name,
                    last_name: formData.state.last_name,
                    email: formData.state.email,
                    address: formData.state.address,
                    area: formData.state.area,
                    postal: formData.state.postal,
                    phone: formData.state.phone,
                },
            }

            dispatch(makePayment(orderData)).then((result) => {
                step.next()
            })
        }
    }

    const handlePrev = () => {
        console.log(formData.state)
        step.prev()
    }

    return (
        <>
            <Review payment={payment}/>
            {/* Show any error or success messages */}
            {message && <Alert severity="error" sx={{
                mb:2
            }}>
                <Typography variant="subtitle1">{message}</Typography>
                </Alert>} 
                <form onSubmit={(e) => handleSubmit(e, elements, stripe)}>
                    <PaymentElement />
                    <Box sx={{
                        display: "flex",
                        justifyContent: "space-between",
                        pt: 3,
                        pb: 3,
                    }}>
                        <BlackButton variant="outlined" onClick={handlePrev} text="Back"/>
                        <BlackButton 
                        type="submit" 
                        disabled={!stripe} 
                        variant="contained" 
                        text={isLoading ? <CircularProgress color="inherit" size={20}/> : `Pay $${formatCurrency(payment.order.subtotal)}`}/>
                    </Box>     
                </form>
        </>
    )
}

export default PaymentForm