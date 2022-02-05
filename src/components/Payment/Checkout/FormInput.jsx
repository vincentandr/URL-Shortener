import React from "react"
import { TextField, Grid } from "@mui/material"
import { useFormContext, Controller } from "react-hook-form"
import MuiPhoneNumber from 'material-ui-phone-number-2';

const FormInput = ({name, placeholder, required, type=""}) => {
    const {control} = useFormContext()

    return (
        <Grid item xs={12} sm={6}>
            <Controller control={control} name={name} render={({field}) => (
                name !== "phone" ? <TextField {...field} fullWidth placeholder={placeholder} required={required} size="small" type={type}/>
                : <MuiPhoneNumber {...field} fullWidth variant="outlined" size="small" defaultCountry={'sg'}/>
            )}/>
        </Grid>
    )
}

export default FormInput