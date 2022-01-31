import React from "react";
import {useState, useContext} from "react";
import {useDispatch} from "react-redux";
import {Box, IconButton, Container, Typography, Button, Table, TableContainer, TableBody, TableRow, TableCell, Drawer, Stack, ButtonGroup} from "@mui/material"
import { createTheme, ThemeProvider } from "@mui/material/styles";
import {DeleteForeverOutlined, Close, Add, Remove} from "@mui/icons-material"

import { removeCartItem, removeAllCartItems, addCartItem } from "../../actions";
import {CartContext} from "../../pages/App"


const Cart = ({drawerState}) => {
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

    const value = useContext(CartContext)
    const dispatch = useDispatch()

    const isEmpty = !value.cart.length

    const handleQty = (op, productId, qty) => {
        // Add +1 to qty if item already exists in cart
        let obj = value.cart.find(item => item.product_id === productId)

        if (op === "increment"){
            obj.qty = obj.qty + qty
        } else if (op === "decrement"){
            obj.qty = obj.qty - qty
        }

        dispatch(addCartItem(productId, obj.qty))
    }
    
    const EmptyCart = () => (
        <Typography variant="h6">
            You have no item in your shopping cart. Shop now!
        </Typography>
    )


    const FilledCart = () => (
        <>
            <TableContainer>
                <Table sx={{ minWidth: 650 }} aria-label="simple table">
                    <TableBody>
                    {value.cart.map((product) => (
                        <TableRow
                        key={product.product_id}
                        sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
                        >
                            <TableCell aling="center">
                                <Box component="img" sx={{
                                    maxHeight: { xs: 167, md: 167 },
                                    maxWidth: { xs: 250, md: 250 },
                                    }}
                                    alt="The house from the offer."
                                    src={product.image}/>
                            </TableCell>
                            <TableCell>
                                <Typography variant="subtitle1">
                                    {product.name}
                                </Typography>
                                <Typography variant="subtitle2">
                                    ${product.price}
                                </Typography>
                                <Typography variant="subtitle2">
                                    <ButtonGroup variant="outlined" size="small" aria-label="outlined primary button group">
                                        <Button onClick={() => handleQty(
                                            "increment",
                                            product.product_id,
                                            1,
                                        )}><Add/></Button>
                                        <ThemeProvider theme={disabledButtonTheme}>
                                            <Button disabled>{product.qty}</Button>
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
                                <div>${product.price * product.qty}</div>
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
                    $1000000
                </Typography>
            </Stack>
            <Stack direction="row" spacing={2}>
                <Button type="button" variant="outlined" color="inherit" fullWidth onClick={() => dispatch(removeAllCartItems())}>Empty Cart</Button>
                <ThemeProvider theme={theme}>
                    <Button type="button" variant="contained" color="primary" fullWidth>Checkout</Button>
                </ThemeProvider>
            </Stack>
        </>
    )

    return(
        <Drawer
            anchor="right"
            open={drawerState}
            onClose={() => value.onClickDrawer(false)}
        >
            <Container>
                <Stack direction="row" spacing={2}>
                    <Typography variant="h4" width="100%">
                        Shopping Cart
                    </Typography>
                    <Button variant="text" size="large" startIcon={<Close />} onClick={() => value.onClickDrawer(false)}>
                        Close
                    </Button>
                </Stack>
                {isEmpty ? <EmptyCart/> : <FilledCart/>}
            </Container>
        </Drawer>
    )
}

export default Cart