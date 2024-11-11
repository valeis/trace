import React, { Component } from "react";
import {Table, Thead, Tbody, Tr, Th, Td, Center} from "@chakra-ui/react";
import axios from "axios";

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
            <Center>
                <Table size="sm" variant="simple" width={"50%"}>
                    <Thead>
                        <Tr>
                            <Th>User ID</Th>
                            <Th>Username</Th>
                            <Th>Email</Th>
                        </Tr>
                    </Thead>
                    <Tbody>
                        {users.map((user) => (
                            <Tr key={user.id}>
                                <Td>{i++}</Td>
                                <Td>{user.username}</Td>
                                <Td>{user.email}</Td>
                            </Tr>
                        ))}
                    </Tbody>
                </Table>
            </Center>
        );
    }
}

export default UserManagement;
