import React from "react";
import {useState} from "react";
import {useDispatch} from "react-redux";
import {useNavigate} from "react-router-dom";
import {Box, Typography, Button, Stack, ButtonGroup, List, ListItem, ListItemText, ListItemButton, Divider} from "@mui/material"
import { createTheme, ThemeProvider } from "@mui/material/styles";
import {DeleteForeverOutlined, Add, Remove} from "@mui/icons-material"

import { removeCartItem, removeAllCartItems, addCartItem, checkout } from "../../../actions";
import { formatCurrency } from "../../../helpers/Utils";

const FilledCart = ({cart}) => {
    const dispatch = useDispatch()
    const navigate = useNavigate()

    const disabledButtonTheme = createTheme({
        palette: {
            action: {
                disabledBackground: 'black',
                disabled: 'inherit',
            }
        }
    });

    const theme = createTheme({
        palette: {
            primary: {
                main: "#212121",
            },
        },
    });

    const handleQty = (op, productId, qty) => {
        // Add +1 to qty if item already exists in cart
        let obj = cart.products.find(item => item.product_id === productId)

        if (obj.qty !== undefined){
            if (op === "increment"){
                qty = obj.qty + qty
            } else if (op === "decrement"){
                qty = obj.qty - qty
            }
        }

        if (qty === 0) {
            dispatch(removeCartItem(productId))
        } else {
            dispatch(addCartItem(productId, qty))
        }
    }

    const handleCheckout = () => {
        dispatch(checkout()).then(() => {
            navigate("/payment")
        })
    }

    return (
        <>
            <List>
                {cart.products.map((product) => (
                    <ListItem key={product.product_id}>
                        <Box component="img" sx={{
                                    maxHeight: { xs: 50, md: 75 },
                                    maxWidth: { xs: 50, md: 75 },
                                    }}
                                    alt="product img"
                                    src={product.image}/>
                        <ListItemText sx={{
                            pl: 2,
                            pr: 2,
                            overflowWrap: "break-word",
                        }}
                        primary={product.name}
                        secondary={`$${formatCurrency(product.price)}`}/>
                        <div>
                            <Typography variant="subtitle1" align="right">
                                ${product.qty === undefined ? 0 : formatCurrency(product.price * product.qty)}
                            </Typography>
                            <ButtonGroup variant="outlined" size="small" aria-label="outlined primary button group">
                                <Button color="inherit" onClick={() => dispatch(removeCartItem(product.product_id))}>
                                    <DeleteForeverOutlined/>
                                </Button>
                                <ThemeProvider theme={disabledButtonTheme}>
                                    <Button disabled>{product.qty === undefined ? 0 : product.qty}</Button>
                                </ThemeProvider>
                                <Button color="inherit" onClick={() => handleQty(
                                    "increment",
                                    product.product_id,
                                    1,
                                )}><Add/></Button>
                                <Button color="inherit" onClick={() => handleQty(
                                    "decrement",
                                    product.product_id,
                                    1,
                                )}><Remove/></Button>
                                
                            </ButtonGroup>
                        </div>
                    </ListItem>
                ))}
                <Divider/>
                <ListItem>
                    <ListItemText primary="Subtotal"/>
                    <Typography variant="subtitle1">
                        ${formatCurrency(cart.subtotal)}
                    </Typography>
                </ListItem>
            </List>
            <Stack direction="row" spacing={2} pb={2}>
                <Button type="button" variant="outlined" color="inherit" fullWidth onClick={() => dispatch(removeAllCartItems())}>Empty Cart</Button>
                <ThemeProvider theme={theme}>
                    <Button type="button" variant="contained" color="primary" fullWidth onClick={handleCheckout}>Checkout</Button>
                </ThemeProvider>
            </Stack>
        </>
    )
}

export default FilledCart