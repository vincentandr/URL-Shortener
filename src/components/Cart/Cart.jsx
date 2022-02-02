import React from "react";
import {useState} from "react";
import {useDispatch} from "react-redux";
import {Box, IconButton, Container, Typography, Button, Table, TableContainer, TableBody, TableRow, TableCell, Drawer, Stack, ButtonGroup} from "@mui/material"
import { createTheme, ThemeProvider } from "@mui/material/styles";
import {DeleteForeverOutlined, Close, Add, Remove} from "@mui/icons-material"

import { removeCartItem, removeAllCartItems, addCartItem, checkout } from "../../actions";
import { formatCurrency } from "../../helpers/Utils";


const Cart = ({cart, drawer}) => {
    const theme = createTheme({
        palette: {
            primary: {
                main: "#212121",
            },
        },
    });

    const disabledButtonTheme = createTheme({
        palette: {
            action: {
                disabledBackground: 'inherit',
                disabled: 'inherit'
            }
        }
    });

    const dispatch = useDispatch()

    const isEmpty = !cart.products.length

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
    
    const EmptyCart = () => (
        <Typography variant="h6">
            You have no item in your shopping cart. Shop now!
        </Typography>
    )


    const FilledCart = () => (
        <>
            <TableContainer>
                <Table aria-label="simple table">
                    <TableBody>
                    {cart.products.map((product) => (
                        <TableRow
                        key={product.product_id}
                        sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
                        >
                            <TableCell aling="center">
                                <Box component="img" sx={{
                                    maxHeight: { xs: 167, md: 167 },
                                    maxWidth: { xs: 250, md: 250 },
                                    }}
                                    alt="product img"
                                    src={product.image}/>
                            </TableCell>
                            <TableCell>
                                <Typography variant="subtitle1">
                                    {product.name}
                                </Typography>
                                <Typography variant="subtitle2">
                                    ${formatCurrency(product.price)}
                                </Typography>
                                <Typography variant="subtitle2">
                                    <ButtonGroup variant="outlined" size="small" aria-label="outlined primary button group">
                                        <Button onClick={() => handleQty(
                                            "increment",
                                            product.product_id,
                                            1,
                                        )}><Add/></Button>
                                        <ThemeProvider theme={disabledButtonTheme}>
                                            <Button disabled>{product.qty === undefined ? 0 : product.qty}</Button>
                                        </ThemeProvider>
                                        <Button onClick={() => handleQty(
                                            "decrement",
                                            product.product_id,
                                            1,
                                        )}><Remove/></Button>
                                    </ButtonGroup>
                                </Typography>
                            </TableCell>
                            <TableCell align="right">
                                <div>${product.qty === undefined ? 0 : formatCurrency(product.price * product.qty)}</div>
                                <IconButton aria-label="Delete item" onClick={() => dispatch(removeCartItem(product.product_id))}>
                                    <DeleteForeverOutlined/>
                                </IconButton>
                            </TableCell>
                        </TableRow>
                    ))}
                    </TableBody>
                </Table>
            </TableContainer>
            <Stack direction="row" pb={2}>
                <Typography variant="h5">
                    Subtotal
                </Typography>
                    <Box sx={{
                        flexGrow: 1
                    }}/>
                <Typography variant="h5">
                    ${cart.subtotal === undefined ? 0 : formatCurrency(cart.subtotal)}
                </Typography>
            </Stack>
            <Stack direction="row" spacing={2} pb={2}>
                <Button type="button" variant="outlined" color="inherit" fullWidth onClick={() => dispatch(removeAllCartItems())}>Empty Cart</Button>
                <ThemeProvider theme={theme}>
                    <Button type="button" variant="contained" color="primary" fullWidth onClick={() => dispatch(checkout())}>Checkout</Button>
                </ThemeProvider>
            </Stack>
        </>
    )

    return(
        <Drawer
            anchor="right"
            open={drawer.state}
            onClose={() => drawer.onClick(false)}
        >
            <Container>
                <Stack direction="row" spacing={2}>
                    <Typography variant="h4" width="100%">
                        Shopping Cart
                    </Typography>
                    <Button variant="text" size="large" startIcon={<Close />} onClick={() => drawer.onClick(false)}>
                        Close
                    </Button>
                </Stack>
                {isEmpty ? <EmptyCart/> : <FilledCart/>}
            </Container>
        </Drawer>
    )
}

export default Cart