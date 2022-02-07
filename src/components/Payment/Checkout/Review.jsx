import React from "react"
import { Box, List, ListItem, ListItemText, Typography, Divider } from "@mui/material"
import { formatCurrency } from "../../../helpers/Utils"

const Review = ({payment}) => {

    return (
        <>
            <Typography variant="h6" >Order Summary</Typography>
            <List sx={{
                maxHeight: "30%",
                overflow: "auto",
            }}>
                {payment.order.items.map((item) => (
                    <ListItem key={item.product_id}>
                        <Box component="img" sx={{
                                    minWidth: {xs: 50, md: 75},
                                    maxHeight: { xs: 50, md: 75 },
                                    maxWidth: { xs: 50, md: 75 },
                                    }}
                                    alt="product img"
                                    src={item.image}/>
                        <ListItemText
                        sx={{
                            pl: 2,
                            pr: 2,
                        }}
                        secondaryTypographyProps={{component: "span"}}
                        primary={item.name} 
                        secondary={<>
                        <div>x{item.qty}</div>
                        <div>{`$${formatCurrency(item.price)}`}</div>
                        </>}/>
                        <Typography variant="body2">
                            ${formatCurrency(item.price * item.qty)}
                        </Typography>
                    </ListItem>
                ))}
            </List>
            <Box sx={{
                display:"flex",
                justifyContent: "space-between",
            }}>
                <Typography variant="h6">Total</Typography>
                <Typography variant="h6">${formatCurrency(payment.order.subtotal)}</Typography>
            </Box>
            <Divider/>
        </>
    )
}

export default Review