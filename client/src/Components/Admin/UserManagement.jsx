import React, { Component } from "react";
import {Table, Thead, Tbody, Tr, Th, Td, Center, Heading, Stack, Flex, Box, Text, Button, Icon} from "@chakra-ui/react";
import axios from "axios";
import { RiAdminFill } from "react-icons/ri";
import { FaUser } from "react-icons/fa";
import {IoIosLogOut} from "react-icons/io";
import { FaChevronCircleRight } from "react-icons/fa";
import {useNavigate} from "react-router-dom";

class UserManagement extends Component {
    state = {
        users: [],
    };

    componentDidMount() {
        this.fetchUsers();
    }

    fetchUsers = async () => {
        const token = localStorage.getItem("access-token").replace(/"/g, ''); // Assume token is stored with key "authToken"
        try {
            const response = await axios.get("http://localhost:1337/user/all/1/10", {
                headers: {
                    Authorization: `${token}`,
                },
            });
            this.setState({ users: response.data.users });
        } catch (error) {
            console.error("Error fetching users:", error);
        }
    };
    render() {
        const { users } = this.state;
        let i = 1;
        return (
            <Stack width="full" gap="2">
                <Flex justifyContent={'flex-end'} marginRight="25%">
                    <Button size="xs" style={{"margin-left":"5px", "margin-top":"3px"}} onClick={() => window.history.back()}>
                        <Icon as={FaChevronCircleRight} w={5} h={5}></Icon>
                    </Button>
                </Flex>
                <Flex justifyContent={'flex-start'} marginLeft="25%" marginBottom="20px">
                    <Heading size="md" style={{"justifyContent":"start"}}>User Management</Heading>
                </Flex>
                <Center>
                    <Table size="sm" variant="simple" width={"50%"}>
                        <Thead>
                            <Tr>
                                <Th>User ID</Th>
                                <Th>Access level</Th>
                                <Th>Username</Th>
                            </Tr>
                        </Thead>
                        <Tbody>
                            {users.map((user) => (
                                <Tr key={user.id}>
                                    <Td>{i++}</Td>
                                    <Td>
                                        <Box
                                            display="flex"
                                            justifyContent="center"
                                        >
                                            {user.role === "admin" && <RiAdminFill />}
                                            {user.role === "user" && <FaUser />}
                                            <Text sx={{ ml: "5px" }}>
                                                {user.role}
                                            </Text>
                                        </Box>
                                    </Td>
                                    <Td>{user.email}</Td>
                                </Tr>
                            ))}
                        </Tbody>
                    </Table>
                </Center>
            </Stack>
        );
    }
}

export default UserManagement;
