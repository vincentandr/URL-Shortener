import React from "react"
import {useContext} from "react";
import { useDispatch } from "react-redux";
import {Card, CardMedia, CardContent, CardActions, Typography, IconButton, Stack, Box} from "@mui/material"
import { AddShoppingCart } from "@mui/icons-material";

import { CartContext } from "../../../pages/App"
import { addCartItem } from "../../../actions";

const Product = ({product}) => {
    const dispatch = useDispatch()
    const value = useContext(CartContext)

    const handleCartClick = (productId, qty) => {
        // Add +1 to qty if item already exists in cart
        let obj = value.cart.find(item => item.product_id === productId)

        if (obj !== undefined){
           qty += obj.qty
        }

        dispatch(addCartItem(productId, qty)).then((result) => {
            value.onClickDrawer(true)
        })
    }

    return (
        <Card>
            <CardMedia image={product.image} title={product.name} style={{height:0, paddingTop: "56.29%"}}/>
            <CardContent>
                <Stack direction="row" spacing={2}>
                    <Typography variant="h6" gutterBottom>
                        {product.name}
                    </Typography>
                    <Box sx={{
                        flexGrow: 1
                    }}/>
                    <Typography variant="h6">
                        ${product.price}
                    </Typography>
                </Stack>
                <Typography variant="body2">{product.desc}</Typography>
            </CardContent>
            <CardActions>
                <Box sx={{
                    flexGrow: 1
                }}/>
                <IconButton aria-label="Add to cart" onClick={() => handleCartClick(
                    product.product_id,
                    1,)}>
                    <AddShoppingCart fontSize="large"/>
                </IconButton>
            </CardActions>
        </Card>
    )
}

export default React.memo(Product);