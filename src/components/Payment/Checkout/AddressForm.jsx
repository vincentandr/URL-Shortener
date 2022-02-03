import React from "react"
import { Typography, Grid, Button } from "@mui/material";
import { useForm, FormProvider } from "react-hook-form";

import FormInput from "./FormInput";

const AddressForm = ({next}) => {
    const methods = useForm()

    return (
        <>
            <Typography variant="h5" gutterBottom>
                Shipping Address
            </Typography>
            <FormProvider {...methods}>
                <form onSubmit={next}>
                    <Grid container spacing={3}>
                        <FormInput required name="firstName" placeholder="First Name"/>
                        <FormInput required name="lastName" placeholder="Last Name"/>
                        <FormInput required name="address" placeholder="Address"/>
                        <FormInput required name="email" placeholder="E-mail" type="email"/>
                        <FormInput required name="area" placeholder="Area"/>
                        <FormInput required name="postal" placeholder="Postal Code"/>
                    </Grid>
                    <Button type="submit" variant="outlined">Next</Button>
                </form>
            </FormProvider>
        </>
    )
}

export default AddressForm;