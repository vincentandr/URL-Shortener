import React from "react"
import { useDispatch } from "react-redux";
import {Card, CardMedia, CardContent, CardActions, Typography, IconButton, Stack, Box} from "@mui/material"
import { AddShoppingCart } from "@mui/icons-material";

import { addCartItem } from "../../../actions";
import { formatCurrency } from "../../../helpers/Utils";

const Product = ({product, cart, onClickDrawer}) => {
    const dispatch = useDispatch()

    const handleCartClick = (productId, qty) => {
        // Add +1 to qty if item already exists in cart
        let obj = cart.products.find(item => item.product_id === productId)

        // If current qty is 0 then skip the if condition
        if (obj !== undefined && obj.qty !== undefined){
            qty += obj.qty
        }

        dispatch(addCartItem(productId, qty)).then((result) => {
            onClickDrawer(true)
        })
    }

    return (
        <Card>
            <CardMedia image={product.image} title={product.name} style={{height:0, paddingTop: "56.29%"}}/>
            <CardContent sx={{
                pb: 0
            }}>
                <Stack direction="row" spacing={2}>
                    <Typography variant="h6" gutterBottom>
                        {product.name}
                    </Typography>
                    <Box sx={{
                        flexGrow: 1
                    }}/>
                    <Typography variant="h6">
                        ${formatCurrency(product.price)}
                    </Typography>
                </Stack>
                <Typography variant="body1">{product.desc}</Typography>
            </CardContent>
            <CardActions>
                <Box sx={{
                    flexGrow: 1,
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