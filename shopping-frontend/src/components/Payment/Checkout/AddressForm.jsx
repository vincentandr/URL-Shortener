import React from "react"
import { Typography, Grid, Stack } from "@mui/material";
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
                        <FormInput required xs={12} sm={6} name="first_name" label="First Name"/>
                        <FormInput required xs={12} sm={6} name="last_name" label="Last Name"/>
                        <FormInput required xs={12} sm={12} name="address" label="Address"/>
                        <FormInput required xs={12} sm={12} name="email" label="E-mail" type="email"/>
                        <FormInput required xs={12} sm={12} name="phone" label="Phone Number"/>
                        <FormInput required xs={12} sm={12} name="area" label="Area"/>
                        <FormInput required xs={12} sm={12} name="postal" label="Postal Code"/>
                    </Grid>
                    <Stack direction="row" justifyContent="space-between" sx={{
                        pt: 3
                    }}>
                        <BlackButton component={Link} to="/" variant="outlined" text="Cancel"/>
                        <BlackButton type="submit" variant="outlined" text="Next"/>
                    </Stack>
                </form>
            </FormProvider>
        </>
    )
}

export default AddressForm;