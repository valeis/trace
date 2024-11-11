import React, { Component } from 'react';
import axios from 'axios';

import SocketConnection from '../../socket-connection';

import {
    Container,
    Flex,
    Textarea,
    Box,
    FormControl,
    FormErrorMessage,
    InputGroup,
    InputRightElement,
    Button,
    Input,
    Tag, Icon, MenuList, MenuButton, MenuItem, Menu, Text
} from '@chakra-ui/react';

import {FaUserCircle, FaUserCog} from "react-icons/fa";
import { FaUserPlus } from "react-icons/fa";
import { IoIosLogOut } from "react-icons/io";
import ChatHistory from './ChatHistory';
import ContactList from './ContactList';
import {Link, Navigate} from "react-router-dom";

class Chat extends Component {
    constructor(props) {
        super(props);
        this.state = {
            socketConn: '',
            username: '',
            message: '',
            to: '',
            isInvalid: false,
            endpoint: 'http://localhost:1337',
            contact: '',
            contacts: [],
            renderContactList: [],
            chats: [],
            chatHistory: [],
            msgs: [],
            redirect: false,
            redirectTo: '/login',
        };
    }

    componentDidMount = async () => {
        const queryParams = new URLSearchParams(window.location.search);
        const user = queryParams.get('u');
        this.setState({ username: user });
        this.getContacts(user);

        const conn = new SocketConnection();
        await this.setState({ socketConn: conn });
        this.state.socketConn.connect(message => {
            const msg = JSON.parse(message.data);

            if (this.state.to === msg.from || this.state.username === msg.from) {
                this.setState(
                    {
                        chats: [...this.state.chats, msg],
                    },
                    () => {
                        this.renderChatHistory(this.state.username, this.state.chats);
                    }
                );
            }
        });

        this.state.socketConn.connected(user);

        console.log('exiting');
    };

    onChange = event => {
        this.setState({ [event.target.name]: event.target.value });
    };

    onSubmit = e => {
        if (e.charCode === 0 && e.code === 'Enter') {
            e.preventDefault();
            const msg = {
                type: 'message',
                chat: {
                    from: this.state.username,
                    to: this.state.to,
                    message: this.state.message,
                },
            };

            this.state.socketConn.sendMsg(msg);
            this.setState({ message: '' });
        }
    };

    getContacts = async user => {
        const res = await axios.get(
            `${this.state.endpoint}/contact-list?username=${user}`
        );
        if (res.data['data'] !== undefined) {
            this.setState({ contacts: res.data.data });
            this.renderContactList(res.data.data);
        }
    };

    fetchChatHistory = async (u1 = 'user1', u2 = 'user2') => {
        const res = await axios.get(
            `${this.state.endpoint}/chat-history?u1=${u1}&u2=${u2}`
        );

        if (res.data['data'] !== undefined) {
            this.setState({ chats: res.data.data.reverse()});
            this.renderChatHistory(u1, res.data.data);
        } else {
            this.setState({ chatHistory: [] });
        }
    };

    addContact = async e => {
        e.preventDefault();
        try {
            const res = await axios.post(`${this.state.endpoint}/verify-contact`, {
                username: this.state.contact,
            });
            if (!res.data.status || !res.data) {
                this.setState({ isInvalid: true });
            } else {
                this.setState({ isInvalid: false });
                let contacts = this.state.contacts;
                contacts.unshift({
                    username: this.state.contact,
                    last_activity: Date.now() / 1000,
                });
                this.renderContactList(contacts);
            }
        } catch (error) {
            console.error(error);
        }
    };

    renderChatHistory = (currentUser, chats) => {
        const history = ChatHistory(currentUser, chats);
        this.setState({ chatHistory: history });
    };

    renderContactList = contacts => {
        const renderContactList = ContactList(contacts, this.sendMessageTo);
        this.setState({ renderContactList });
    };

    sendMessageTo = to => {
        this.setState({ to });
        this.fetchChatHistory(this.state.username, to);
    };

    logout = () => {
        localStorage.removeItem("access-token");
        localStorage.removeItem("refresh-token");
        localStorage.removeItem("username");
        localStorage.removeItem("userId");
        this.setState({redirect: true});
    }

    render() {
        return (
            <Container>
                {this.state.redirect && (
                    <Navigate to={this.state.redirectTo} replace={true}></Navigate>
                )}
                <Flex justifyContent={'flex-end'}>
                    <Menu>
                        <MenuButton as={Tag} cursor="pointer" padding="5px 10px" borderRadius="full" bg="gray.100" _hover={{ bg: "gray.200" }}>
                            <Flex alignItems="center">
                                <Box marginRight="2">
                                    <FaUserCircle color="#4A5568" size="1.5em" />
                                </Box>
                                    {this.state.username}
                            </Flex>
                        </MenuButton>
                        <MenuList border="1px solid" borderColor="gray.200" boxShadow="lg" borderRadius="md" bg="white">
                            <MenuItem as={Link} to="/users" icon={<FaUserCog />} _hover={{ bg: "gray.100" }}>
                                <Text fontWeight="small">User Management</Text>
                            </MenuItem>
                        </MenuList>
                    </Menu>

                    <Button size="xs" style={{"margin-left":"5px", "margin-top":"3px"}} onClick={this.logout}>
                        <Icon as={IoIosLogOut} w={5} h={5}></Icon>
                    </Button>
                </Flex>
                <Box p="10px"/>
                <Container paddingBottom={2}>
                    <Box>
                        <FormControl isInvalid={this.state.isInvalid}>
                            <InputGroup size="sm">
                                <Input
                                    variant="flushed"
                                    type="text"
                                    placeholder="Add Contact"
                                    name="contact"
                                    value={this.state.contact}
                                    onChange={this.onChange}
                                    style={{borderBottom: !this.state.isInvalid ? '0' : '', boxShadow: !this.state.isInvalid ? 'none' : ''}}
                                />
                                <InputRightElement width="6rem">
                                    <Button
                                        colorScheme={'blue'}
                                        h="2rem"
                                        size="sm"
                                        variant="solid"
                                        type="submit"
                                        onClick={this.addContact}
                                        marginLeft={'25px'}
                                    >
                                        <Box
                                            marginRight={'5px'}>
                                            <FaUserPlus/>
                                        </Box>
                                        Add
                                    </Button>
                                </InputRightElement>
                            </InputGroup>
                            { !this.state.isInvalid ? null : (
                                <FormErrorMessage>Contact does not exist</FormErrorMessage>
                            )}
                        </FormControl>
                    </Box>
                </Container>
                <Flex>
                    <Box
                        textAlign={'left'}
                        overflowY={'scroll'}
                        flex="1"
                        h={'32rem'}
                        borderWidth={1}
                        borderRightWidth={0}
                        borderRadius={'xl'}
                    >
                        {this.state.renderContactList}
                    </Box>

                    <Box flex="2">
                        <Container
                            borderWidth={1}
                            borderLeftWidth={0}
                            borderBottomWidth={0}
                            borderRadius={'xl'}
                            textAlign={'right'}
                            h={'28rem'}
                            padding={2}
                            overflowY={'scroll'}
                            display="flex"
                            flexDirection={'column-reverse'}
                        >
                            {this.state.chatHistory}
                        </Container>

                        <Box flex="1">
                            <FormControl onKeyDown={this.onSubmit} onSubmit={this.onSubmit}>
                                <Textarea
                                    fontSize={'sm'}
                                    type="submit"
                                    borderWidth={1}
                                    borderRadius={'xl'}
                                    minH={'4rem'}
                                    placeholder="Type a message"
                                    size="lg"
                                    resize={'none'}
                                    name="message"
                                    value={this.state.message}
                                    onChange={this.onChange}
                                    isDisabled={this.state.to === ''}
                                />
                            </FormControl>
                        </Box>
                    </Box>
                </Flex>
            </Container>
        );
    }
}

export default Chat;