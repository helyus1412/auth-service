export interface User {
  id: number;
  name: string;
  email: string;
}

let users: User[] = [
  { id: 1, name: 'John Doe', email: 'john@example.com' },
  { id: 2, name: 'Jane Smith', email: 'jane@example.com' },
];

export const getUsers = () => users;
export const getUserById = (id: number) => users.find(u => u.id === id);
export const addUser = (user: User) => { users.push(user); };
export const updateUser = (id: number, updated: Partial<User>) => {
  users = users.map(u => (u.id === id ? { ...u, ...updated } : u));
};
export const deleteUser = (id: number) => {
  users = users.filter(u => u.id !== id);
};
