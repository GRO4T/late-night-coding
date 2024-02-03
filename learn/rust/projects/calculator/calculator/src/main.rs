use std::io;
use std::io::Write;

fn main() {
    loop {
        let mut expr = String::new();

        print!(">> ");
        io::stdout().flush().unwrap();
        
        io::stdin()
            .read_line(&mut expr)
            .expect("Failed to read expression");

        if expr == "q\n" {
            break;
        }

        // first handle the simplest case
        expr = expr.replace(" ", "");
        let left_op = expr.chars().nth(0).unwrap() as i32 - '0' as i32;
        let op = expr.chars().nth(1).unwrap();
        let right_op = expr.chars().nth(2).unwrap() as i32 - '0' as i32;

        let result_or_none = match op {
            '+' => Some(left_op + right_op),
            '-' => Some(left_op - right_op),
            '*' => Some(left_op * right_op),
            '/' => Some(left_op / right_op),
            _ => None::<i32>,
        };
        
        if result_or_none.is_none() {
            println!("Unknown expression!");
        }
        else {
            let result = result_or_none.unwrap();
            println!("= {result}");
        }
    }
}
