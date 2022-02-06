import React from "react"
import { TextField, Grid } from "@mui/material"
import { useFormContext, Controller } from "react-hook-form"
import MuiPhoneNumber from 'material-ui-phone-number-2';

const FormInput = ({name, label, xs, sm, required, type}) => {
    const {control} = useFormContext()

    return (
        <Grid item xs={xs} sm={sm}>
            <Controller control={control} name={name} render={({field}) => (
                name !== "phone" ? <TextField {...field} label={label} fullWidth required={required} size="small" type={type}/>
                : <MuiPhoneNumber {...field} fullWidth variant="outlined" label={label} size="small" defaultCountry={'sg'}/>
            )}/>
        </Grid>
    )
}

export default FormInput