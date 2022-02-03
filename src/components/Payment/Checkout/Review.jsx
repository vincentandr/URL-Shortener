import React from "react"
import { List, ListItem, ListItemText, Typography } from "@mui/material"
import { formatCurrency } from "../../../helpers/Utils"

const Review = ({cart}) => {
    return (
        <>
            <Typography variant="h6" gutterBottom>Order Summary</Typography>
            <List disablePadding>
                {cart.products.map((item) => (
                    <ListItem key={item.product_id}>
                        <ListItemText primary={item.name} secondary={item.qty}/>
                        <Typography variant="body2">
                            ${formatCurrency(item.price * item.qty)}
                        </Typography>
                    </ListItem>
                ))}
                <ListItem>
                    <ListItemText primary="Total"/>
                    <Typography variant="subtitle1">${formatCurrency(cart.subtotal)}</Typography>
                </ListItem>
            </List>
        </>
    )
}

export default Review