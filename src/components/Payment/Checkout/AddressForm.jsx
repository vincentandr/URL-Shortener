import React from "react"
import { Box, Typography, Grid } from "@mui/material";
import { useForm, FormProvider } from "react-hook-form";
import {Link} from "react-router-dom"
import { BlackButton } from "../../../theme";

import FormInput from "./FormInput";

const AddressForm = ({next, formData}) => {
    const methods = useForm({
        defaultValues: {
            first_name: formData.state.first_name,
            last_name: formData.state.last_name,
            address: formData.state.address,
            email: formData.state.email,
            area: formData.state.area,
            postal: formData.state.postal,
            phone: formData.state.phone,
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
                        <FormInput required name="first_name" placeholder="First Name"/>
                        <FormInput required name="last_name" placeholder="Last Name"/>
                        <FormInput required name="address" placeholder="Address"/>
                        <FormInput required name="email" placeholder="E-mail" type="email"/>
                        <FormInput required name="area" placeholder="Area"/>
                        <FormInput required name="postal" placeholder="Postal Code"/>
                        <FormInput required name="phone" placeholder="Phone Number"/>
                    </Grid>
                    <Box sx={{
                        display: "flex",
                        justifyContent: "space-between",
                        pt: 3,
                        pb: 3,
                    }}>
                        <BlackButton component={Link} to="/" variant="outlined" text="Cancel"/>
                        <BlackButton type="submit" variant="outlined" text="Next"/>
                    </Box>
                </form>
            </FormProvider>
        </>
    )
}

export default AddressForm;