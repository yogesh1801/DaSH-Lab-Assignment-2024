#ifndef LOCKS_CPP
#define LOCKS_CPP

#include <string>
#include <stdexcept>

class Lock {
public:
    std::string type;
    std::string table_name;
    int row;
    int column;

    Lock(std::string lock_type, std::string table, int r = -1, int c = -1)
        : type(lock_type), table_name(table), row(r), column(c) 
    {
        if (type == "TABLE") {
            if (table_name.empty()) {
                throw std::invalid_argument("Table name is required for TABLE lock.");
            }
            row = -1;
            column = -1;
        } 
        else if (type == "ROW") {
            if (table_name.empty() || row == -1) {
                throw std::invalid_argument("Table name and row are required for ROW lock.");
            }
            column = -1;
        } 
        else if (type == "COLUMN") {
            if (table_name.empty() || column == -1) {
                throw std::invalid_argument("Table name and column are required for COLUMN lock.");
            }
            row = -1;
        } 
        else if (type == "CELL") {
            if (table_name.empty() || row == -1 || column == -1) {
                throw std::invalid_argument("Table name, row, and column are required for CELL lock.");
            }
        } 
        else {
            throw std::invalid_argument("Invalid lock type. Allowed types are TABLE, ROW, COLUMN, CELL.");
        }
    }
};

#endif
