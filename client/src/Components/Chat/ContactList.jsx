import React from 'react';

import {Text, Box, Divider, Flex} from '@chakra-ui/react';
import {formatDate} from "../../shared/date-time/data-time-pipe";

const ContactList = (contacts, sendMessage) => {
    const contactList = contacts.map(c => {
        const ts = new Date(c.last_activity * 1000);

        return (
            <div key={c.username}>
                <Box
                    as="button"
                    textAlign={'left'}
                    key={c.username}
                    p={2}
                    marginTop={2}
                    marginBottom={2}
                    paddingRight={2}
                    borderRadius={20}
                    borderColor="-moz-initial"
                    borderBottomColor={'whiteAlpha.500'}
                    onClick={() => sendMessage(c.username)}
                    width={'100%'}
                >
                    <Text fontSize={'md'}>{c.username}</Text>
                    <Flex
                        justify={'flex-end'}
                    >
                        <Text as={'sub'} fontSize="xs">
                            {' '}
                            {formatDate(ts)}{' '}
                        </Text>
                    </Flex>
                </Box>
                <Divider></Divider>
            </div>
        );
    });

    return contactList;
};

export default ContactList;