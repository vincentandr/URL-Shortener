import React from "react"
import { Box, Typography, Grid, Button } from "@mui/material";
import { useForm, FormProvider } from "react-hook-form";
import {Link} from "react-router-dom"

import FormInput from "./FormInput";

const AddressForm = ({next, formData}) => {
    const methods = useForm({
        defaultValues: {
            firstName: formData.state.firstName,
            lastName: formData.state.lastName,
            address: formData.state.address,
            email: formData.state.email,
            area: formData.state.area,
            postal: formData.state.postal,
        }
    })

    return (
        <>
            <Typography variant="h6" gutterBottom>
                Shipping Address
            </Typography>
            <FormProvider {...methods}>
                <form onSubmit={methods.handleSubmit((data) => {
                    formData.set(data)
                    next();
                })}>
                    <Grid container spacing={3}>
                        <FormInput required name="firstName" placeholder="First Name"/>
                        <FormInput required name="lastName" placeholder="Last Name"/>
                        <FormInput required name="address" placeholder="Address"/>
                        <FormInput required name="email" placeholder="E-mail" type="email"/>
                        <FormInput required name="area" placeholder="Area"/>
                        <FormInput required name="postal" placeholder="Postal Code"/>
                    </Grid>
                    <Box sx={{
                        display: "flex",
                        justifyContent: "space-between",
                        pt: 3,
                        pb: 3,
                    }}>
                        <Button component={Link} to="/" variant="outlined">Cancel</Button>
                        <Button type="submit" variant="outlined">Next</Button>
                    </Box>
                </form>
            </FormProvider>
        </>
    )
}

export default AddressForm;