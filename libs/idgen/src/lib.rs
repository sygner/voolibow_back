use rand::{Rng, distributions::Alphanumeric};

pub fn alpha_numeric(len: usize) -> String{
    rand::prelude::thread_rng()
    .sample_iter(Alphanumeric)
    .take(len)
    .map(char::from)
    .collect::<String>()
}

pub fn numeric_code_i32(min:i32,max:i32)->i32{
    rand::prelude::thread_rng().gen_range(min..max)
}

pub fn numeric_code_i16(min:i16,max:i16)->i16{
    rand::prelude::thread_rng().gen_range(min..max)
}

pub fn numeric_code_u8(min:u8,max:u8)->u8{
    rand::prelude::thread_rng().gen_range(min..max)
}

pub fn numeric_code_u32(min:u32,max:u32)->u32{
    rand::prelude::thread_rng().gen_range(min..max)
}

pub fn numeric_code_u64(min:u64,max:u64)->u64{
    rand::prelude::thread_rng().gen_range(min..max)
}

pub fn numeric_code_usize(min:usize,max:usize)->usize{
    rand::prelude::thread_rng().gen_range(min..max)
}

pub fn license_alpha_numeric()->String{
    let id =  rand::prelude::thread_rng()
    .sample_iter(Alphanumeric)
    .take(16)
    .map(char::from)
    .collect::<String>();
    let mut reg = "".to_owned();
    for i in id.chars().into_iter().enumerate(){
        if i.0 !=0{
            if i.0 % 4 == 0 {
                reg.push('-')
            }
        }
        reg.push(i.1);
    }
    reg
}
