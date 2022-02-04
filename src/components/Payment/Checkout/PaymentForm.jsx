import React from "react"
import { useDispatch } from "react-redux"
import { Box, Button, Typography } from "@mui/material"
import { Elements, CardElement, ElementsConsumer } from "@stripe/react-stripe-js"
import { loadStripe } from "@stripe/stripe-js"

import Review from "./Review"
import { formatCurrency } from "../../../helpers/Utils"
import { makePayment } from "../../../actions"

const stripePromise = loadStripe(process.env.REACT_APP_STRIPE_PUBLIC_KEY)

const PaymentForm = ({cart, step, formDataState}) => {
    const dispatch = useDispatch()

    const handleSubmit = async (event, elements, stripe) => {
        event.preventDefault();

        if (!stripe || !elements) return;

        const cardElement = elements.getElement(CardElement);

        const {error, paymentMethod} = await stripe.createPaymentMethod({type: "card", card: cardElement})

        if(error) {
            console.log(error)
        }else {
            const orderData = {
                order_id: cart.order_id,
                customer: {
                    firstname: formDataState.firstName,
                    lastname: formDataState.lastname,
                    email: formDataState.email,
                    address: formDataState.address,
                    city: formDataState.city},
                paymentMethod: {gateway: "stripe", stripe: {payment_method_id: paymentMethod.id}},
            }

            dispatch(makePayment(orderData)).then((result) => {
                step.next()
            })
        }
    }

    const handlePrev = () => {
        console.log(formDataState)
        step.prev()
    }

    return (
        <>
            <Review cart={cart}/>
            <Typography variant="h6" gutterBottom>Payment method</Typography>
            <Elements stripe={stripePromise}>
                <ElementsConsumer>
                    {({elements, stripe}) => (
                        <form onSubmit={(e) => handleSubmit(e, elements, stripe)}>
                            <CardElement/>
                            <Box sx={{
                                display: "flex",
                                justifyContent: "space-between",
                                pt: 3,
                                pb: 3,
                            }}>
                                <Button variant="outlined" onClick={handlePrev}>Back</Button>
                                <Button type="submit" disabled={!stripe} variant="contained">Pay ${formatCurrency(cart.subtotal)}</Button>  
                            </Box>      
                        </form>
                    )}
                </ElementsConsumer>
            </Elements>
        </>
    )
}

export default PaymentForm