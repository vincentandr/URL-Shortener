import React from "react"
import { TextField, Grid } from "@mui/material"
import { useFormContext, Controller } from "react-hook-form"

const FormInput = ({name, placeholder, required, type=""}) => {
    const {control} = useFormContext()

    return (
        <Grid item xs={12} sm={6}>
            <Controller control={control} name={name} render={({field}) => (
                <TextField {...field} fullWidth placeholder={placeholder} required={required} size="small" type={type}/>
            )}/>
        </Grid>
    )
}

export default FormInput