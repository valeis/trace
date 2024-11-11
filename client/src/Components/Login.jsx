import React, {Component} from 'react';
import axios from 'axios';

import {
    Box,
    Button,
    Container,
    Flex,
    FormControl,
    FormErrorMessage,
    FormLabel,
    Heading,
    Input,
    Stack,
    Text,
    useColorModeValue
} from '@chakra-ui/react';

import {Navigate} from 'react-router-dom';
import {EditIcon} from '@chakra-ui/icons';
import {Blur} from "./Register";
import Chat from "./Chat/Chat";

class Login extends Component {
    constructor(props) {
        super(props);
        this.state = {
            username: '',
            password: '',
            message: '',
            isInvalid: false,
            endpoint: 'http://localhost:1337/login',
            redirect: false,
            redirectTo: !localStorage.getItem("username") ? '/chat?u=' : `/chat?u=${localStorage.getItem("username")}`,
        };
    }

    componentDidMount() {
        const token = localStorage.getItem("access-token");
        if (token) {
            this.setState({ redirect: true})
        }
    }

    onChange = event => {
        this.setState({[event.target.name]: event.target.value});
    };

    onSubmit = async e => {
        e.preventDefault();

        try {
            const res = await axios.post(this.state.endpoint, {
                email: this.state.username,
                password: this.state.password,
            });

            if (res.data.status) {
                const redirectTo = this.state.redirectTo + this.state.username;
                localStorage.setItem("access-token", JSON.stringify(res.data["access-token"]));
                localStorage.setItem("refresh-token", JSON.stringify(res.data["refresh-token"]));
                localStorage.setItem("userId", JSON.stringify(res.data["userId"]));
                localStorage.setItem("username", this.state.username)
                this.setState({redirect: true, redirectTo});
            } else {
                this.setState({message: res.data.message, isInvalid: true});
            }
        } catch (error) {
            this.setState({message: 'something went wrong', isInvalid: true});
        }
    };

    render() {
        return (
            <Flex
                align="center"
                justify="center"
                minH="100vh"
            >
                {this.state.redirect && (
                    <Navigate to={this.state.redirectTo} replace={true}></Navigate>
                )}
                <Container
                    textAlign="center"
                >
                        <Flex
                            justify="center"
                            align="center"
                            minW="2xl"
                        >
                            <Stack
                                bg={'gray.50'}
                                rounded={'xl'}
                                p={{base: 4, sm: 6, md: 8}}
                                spacing={{base: 8}}
                                // maxW={{lg: 'lg'}}
                                width="70%"
                            >
                                <Stack spacing={4}>
                                    <Flex justifyContent="center">
                                        <Heading
                                            color={'gray.800'}
                                            lineHeight={1.1}
                                            fontSize={{base: '2xl', sm: '3xl', md: '4xl'}}>
                                            Login
                                        </Heading>
                                    </Flex>
                                </Stack>
                                <FormControl isInvalid={this.state.isInvalid}>
                                    <FormLabel color={'gray.800'}>Username</FormLabel>
                                    <Input
                                        type="text"
                                        placeholder="Username"
                                        name="username"
                                        bg={'gray.100'}
                                        border={0}
                                        color={'gray.500'}
                                        _placeholder={{
                                            color: 'gray.500',
                                        }}
                                        value={this.state.username}
                                        onChange={this.onChange}
                                    />
                                </FormControl>
                                <FormControl isInvalid={this.state.isInvalid}>
                                    <FormLabel color={'gray.800'}>Password</FormLabel>
                                    <Input
                                        type="password"
                                        name="password"
                                        value={this.state.password}
                                        placeholder="Password"
                                        bg={'gray.100'}
                                        border={0}
                                        color={'gray.500'}
                                        _placeholder={{
                                            color: 'gray.500',
                                        }}
                                        onChange={this.onChange}
                                    />
                                    {!this.state.isInvalid ? (
                                        ''
                                    ) : (
                                        <FormErrorMessage>
                                            invalid username or password
                                        </FormErrorMessage>
                                    )}
                                </FormControl>
                                <Button
                                    size="lg"
                                    leftIcon={<EditIcon/>}
                                    colorScheme="green"
                                    variant="solid"
                                    type="submit"
                                    fontFamily={'heading'}
                                    mt={2}
                                    w={'full'}
                                    bgGradient="linear(to-r, red.400,pink.400)"
                                    color={'white'}
                                    onClick={this.onSubmit}
                                    _hover={{
                                        bgGradient: 'linear(to-r, red.400,pink.400)',
                                        boxShadow: 'xl',
                                    }}
                                >
                                    Login
                                </Button>
                            </Stack>
                        </Flex>
                </Container>
                <Blur position={'absolute'} top={-10} left={-10} style={{filter: 'blur(70px)'}}/>
            </Flex>
        );
    }
}

export default Login;