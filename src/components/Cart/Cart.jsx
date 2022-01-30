import React from "react";
import {useState, useContext} from "react";
import {useDispatch} from "react-redux";
import {Box, IconButton, Container, Typography, Button, Table, TableContainer, TableBody, TableRow, TableCell, Drawer, Paper, Stack} from "@mui/material"
import {DeleteForeverOutlined, Close} from "@mui/icons-material"

import { removeCartItem, removeAllCartItems } from "../../actions";
import {CartContext} from "../../pages/App"

const Cart = ({drawerState}) => {
    const value = useContext(CartContext)
    const dispatch = useDispatch()

    const isEmpty = !value.cart.length
    
    const EmptyCart = () => (
        <Typography variant="h6">
            You have no item in your shopping cart. Shop now!
        </Typography>
    )

    const FilledCart = () => (
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
                            height: 233,
                            width: 350,
                            maxHeight: { xs: 233, md: 167 },
                            maxWidth: { xs: 350, md: 250 },
                            }}
                            alt="The house from the offer."
                            src={product.image}/>
                    </TableCell>
                    <TableCell align="right">{product.name}</TableCell>
                    <TableCell align="right">{product.qty}</TableCell>
                    <TableCell align="right">{product.price}</TableCell>
                    <TableCell align="right">{product.desc}</TableCell>
                    <TableCell align="right">
                        <IconButton aria-label="Delete item" onClick={() => dispatch(removeCartItem(product.product_id))}>
                            <DeleteForeverOutlined fontSize="large"/>
                        </IconButton>
                    </TableCell>
                    </TableRow>
                ))}
                </TableBody>
            </Table>
        </TableContainer>
    )

    return(
        <Drawer
            anchor="right"
            open={drawerState}
            onClose={() => value.onClickDrawer(false)}
        >
            <Container>
                <div>
                    <Typography variant="h4" display="inline">
                        Shopping Cart
                    </Typography>
                    <Button variant="text" size="large" startIcon={<Close />} onClick={() => value.onClickDrawer(false)}>
                        Close
                    </Button>
                </div>
                {isEmpty ? <EmptyCart/> : <FilledCart/>}
                <div>
                    <Button type="button" variant="contained" color="secondary" onClick={() => dispatch(removeAllCartItems())}>Empty Cart</Button>
                    <Button type="button" variant="contained" color="primary">Checkout</Button>
                </div>
            </Container>
        </Drawer>
    )
}

export default Cart