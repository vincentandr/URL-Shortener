import React from "react";
import {Dialog, DialogTitle, DialogContent, DialogActions, TextField, Stack, Button} from "@mui/material"


const Login = ({loginState, onClickLogin}) => {
    const handleSubmit = () => {
        
    }

    return (
        <form onSubmit={handleSubmit}>
            <Dialog 
                disablePortal
                open={loginState}
                onClose={() => onClickLogin(false)}
                fullWidth={true}
                maxWidth="xs"
            >
                <DialogTitle>
                    Login
                </DialogTitle>
                <DialogContent>
                    <Stack spacing={2}>
                        <TextField
                            required
                            variant="outlined"
                            id="contained-required"
                            defaultValue=""
                            placeholder="Email"
                        />
                        <TextField
                            required
                            variant="outlined"
                            id="outlined-disabled"
                            defaultValue=""
                            placeholder="Password"
                        />
                    </Stack>
                </DialogContent>
                <DialogActions>
                    <Button type="submit" variant="outlined" fullWidth={true}>submit</Button>
                </DialogActions>
            </Dialog>
        </form>
    )
}

export default Login;