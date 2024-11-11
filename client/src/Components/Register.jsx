import React, { Component } from 'react';
import axios from 'axios';

import {
    Container,
    FormControl,
    FormLabel,
    FormErrorMessage,
    Box,
    Input,
    Stack,
    Button,
    IconProps,
    Heading,
    Text,
    SimpleGrid,
    useBreakpointValue,
    Icon,
} from '@chakra-ui/react';

import { Navigate } from 'react-router-dom';
export const Blur = (props: IconProps) => {
    return (
        <Icon
            width={useBreakpointValue({ base: '100%', md:"40vw", lg: '30vw'})}
            zIndex={useBreakpointValue({ base: -1, md: -1, lg: 0})}
            height="560px"
            viewBox="0 0 528 560"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
            {...props}>
            <circle cx="71" cy="61" r="111" fill="#F56565" />
            <circle cx="244" cy="106" r="139" fill="#ED64A6" />
            <circle cy="291" r="139" fill="#ED64A6" />
            <circle cx="80.5" cy="189.5" r="101.5" fill="#ED8936" />
            <circle cx="196.5" cy="317.5" r="101.5" fill="#ECC94B" />
            <circle cx="70.5" cy="458.5" r="101.5" fill="#48BB78" />
            <circle cx="426.5" cy="-0.5" r="101.5" fill="#4299E1" />
        >
        </Icon>
    )
}

class Register extends Component {
    constructor(props) {
        super(props);
        this.state = {
            username: '',
            password: '',
            message: '',
            isInvalid: '',
            endpoint: 'http://localhost:1337/register',
            redirect: false,
            redirectTo: '/login',
        };
    }

    componentDidMount() {
        const token = localStorage.getItem("access-token");
        if (token) {
            this.setState({redirect: true})
        }
    }
    onChange = event => {
        this.setState({ [event.target.name]: event.target.value });
    };

    onSubmit = async e => {
        e.preventDefault();

        try {
            const res = await axios.post(this.state.endpoint, {
                email: this.state.username,
                password: this.state.password,
            });

            if (res.data.status) {
                const redirectTo = this.state.redirectTo;
                this.setState({ redirect: true, redirectTo });
            } else {
                // on failed
                this.setState({ message: res.data.message, isInvalid: true });
            }
        } catch (error) {
            console.log(error);
            this.setState({ message: 'something went wrong', isInvalid: true });
        }
    };
    render() {
        if (this.state.redirect) {
            return <Navigate to={this.state.redirectTo}/>
        }

        return (
            <div>
                <Box position={'relative'}>
                    <Container
                        as={SimpleGrid}
                        maxW={'7xl'}
                        columns={{ base: 1, md: 2}}
                        spacing={{ base:10, lg: 32}}
                        py={{ base:10, sm:20, lg:32}}>
                        <Stack spacing={{base:10, md:20}}>
                            <Box paddingY="2" />
                            <Heading
                                lineHeight={1.1}
                                fontSize={{ base: '3xl', sm: '4xl', md: '5xl', lg: '6xl'}}>
                                Connect with friends{' '}
                                <Text as={'span'} bgGradient="linear(to-r, red.400,pink.400)" bgClip="text">
                                    &
                                </Text>{' '}
                                family instantly!
                            </Heading>
                        </Stack>
                        <Stack
                            bg={'gray.50'}
                            rounded={'xl'}
                            p={{ base: 4, sm: 6, md: 8 }}
                            spacing={{ base: 8 }}
                            maxW={{ lg: 'lg' }}>
                            <Stack spacing={4}>
                                <Heading
                                    color={'gray.800'}
                                    lineHeight={1.1}
                                    fontSize={{ base: '2xl', sm: '3xl', md: '4xl' }}>
                                    Create Your Account
                                </Heading>
                            </Stack>
                            <Box as={'form'} mt={10}>
                                <Stack spacing={4}>
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
                                        {!this.state.isInvalid ? (
                                            <></>
                                        ) : (
                                            <FormErrorMessage>{this.state.message}</FormErrorMessage>
                                        )}
                                    </FormControl>
                                    <FormControl>
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
                                    </FormControl>
                                </Stack>
                                <Button
                                    fontFamily={'heading'}
                                    mt={8}
                                    w={'full'}
                                    bgGradient="linear(to-r, red.400,pink.400)"
                                    color={'white'}
                                    onClick={this.onSubmit}
                                    _hover={{
                                        bgGradient: 'linear(to-r, red.400,pink.400)',
                                        boxShadow: 'xl',
                                    }}
                                >
                                    Register
                                </Button>
                            </Box>
                            form
                        </Stack>
                    </Container>
                    <Blur position={'absolute'} top={-10} left={-10} style={{ filter: 'blur(70px)' }} />
                </Box>
            </div>
        );
    }
}

export default Register;